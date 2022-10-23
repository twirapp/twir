import Twitch from '@nestjs-hybrid-auth/twitch';
import {
  BadRequestException,
  CacheTTL,
  CACHE_MANAGER,
  Controller,
  ExecutionContext,
  Get,
  HttpException,
  HttpStatus,
  Inject,
  Injectable,
  Post,
  Query,
  Req,
  Res,
  UseGuards,
  UseInterceptors,
} from '@nestjs/common';
import { AuthGuard, IAuthModuleOptions } from '@nestjs/passport';
import CacheManager from 'cache-manager';
import Express from 'express';
import merge from 'lodash.merge';

import { CustomCacheInterceptor } from '../helpers/customCacheInterceptor.js';
import { JwtAuthGuard } from '../jwt/jwt.guard.js';
import { JwtAuthService } from '../jwt/jwt.service.js';
import { scope } from './auth.module.js';
import { AuthService } from './auth.service.js';

const REFRESH_TOKEN = 'refresh_token';
const REFRESH_TOKEN_EXPIRE_TIME = 1000 * 60 * 60 * 24 * 31; // 31day

@Injectable()
class TwitchAuthGuard extends AuthGuard('twitch') {
  constructor(options?: Twitch.TwitchAuthGuardOptions) {
    super(
      merge({}, options, {
        property: 'hybridAuthResult',
      }),
    );
  }

  getAuthenticateOptions(
    context: ExecutionContext,
  ): Promise<IAuthModuleOptions> | IAuthModuleOptions | undefined {
    const req = context.switchToHttp().getRequest() as Express.Request;

    return req.query;
  }
}

function UseTwitchAuth(options?: Twitch.TwitchAuthGuardOptions) {
  return UseGuards(new TwitchAuthGuard(options));
}

@Controller('auth')
export class AuthController {
  constructor(
    private readonly jwtAuthService: JwtAuthService,
    private readonly authService: AuthService,
    @Inject(CACHE_MANAGER) private readonly cacheManager: CacheManager.Cache,
  ) {}

  @UseTwitchAuth()
  @Get()
  login() {
    return 'Login with twitch...';
  }

  @Get('token')
  async twitchCallback(
    @Query() query: { code: string; state: string },
    @Res() res: Express.Response,
  ) {
    try {
      const jwtPayload = await this.authService.authorizeUserByTwitch(query.code, query.state);

      const { accessToken, refreshToken } = await this.jwtAuthService.generateKeypair(jwtPayload);

      res
        .cookie(REFRESH_TOKEN, refreshToken, { httpOnly: true, maxAge: REFRESH_TOKEN_EXPIRE_TIME })
        .send({ accessToken });
    } catch (error) {
      res.status(400).send('Something wrong with your authorization. Please try authorize again.');
    }
  }

  @Post('token')
  async refresh(@Req() req: Express.Request, @Res() res: Express.Response) {
    const refreshToken = req.cookies[REFRESH_TOKEN] as string | undefined;
    if (!refreshToken) throw new BadRequestException('Refresh token not passed to body');

    try {
      const accessToken = await this.jwtAuthService.refreshAccessToken(refreshToken);
      res.send({ accessToken });
    } catch (error) {
      res.status(400).send('Something wrong with your authorization. Please try authorize again.');
    }
  }

  @Post('logout')
  @UseGuards(JwtAuthGuard)
  async logout(@Req() req: Express.Request, @Res() res: Express.Response) {
    await this.cacheManager.del(`nest:cache:auth/profile:${req.user.id}`);

    res.clearCookie(REFRESH_TOKEN).send(true);
  }

  @CacheTTL(600)
  @UseInterceptors(
    CustomCacheInterceptor((ctx) => {
      const req = ctx.switchToHttp().getRequest() as Express.Request;
      return `nest:cache:auth/profile:${req.user.id}`;
    }),
  )
  @UseGuards(JwtAuthGuard)
  @Get('profile')
  async showProfile(@Req() req: Express.Request) {
    const isHasNeededScopes = req.user.scopes
      ? scope.every((s) => req.user.scopes.includes(s))
      : false;

    if (!isHasNeededScopes) {
      throw new HttpException('Missed scopes', HttpStatus.UNAUTHORIZED);
    }

    return await this.authService.getProfile(req.user);
  }
}
