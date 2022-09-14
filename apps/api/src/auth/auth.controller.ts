import Twitch from '@nestjs-hybrid-auth/twitch';
import {
  BadRequestException,
  Body,
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
import { config } from '@tsuwari/config';
import { exchangeCode, getTokenInfo } from '@twurple/auth';
import CacheManager from 'cache-manager';
import Express from 'express';
import merge from 'lodash.merge';

import { CustomCacheInterceptor } from '../helpers/customCacheInterceptor.js';
import { JwtAuthGuard } from '../jwt/jwt.guard.js';
import { JwtAuthService } from '../jwt/jwt.service.js';
import { scope } from './auth.module.js';
import { AuthService } from './auth.service.js';

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
  @Get('')
  login() {
    return 'Login with twitch...';
  }

  @Get('token')
  async callback(@Res() res: Express.Response, @Query() query: { code: string; state: string }) {
    const code = await exchangeCode(
      config.TWITCH_CLIENTID,
      config.TWITCH_CLIENTSECRET,
      query.code,
      Buffer.from(query.state, 'base64').toString('utf-8'),
    );
    const tokenInfo = await getTokenInfo(code.accessToken, config.TWITCH_CLIENTID);

    const { accessToken, refreshToken } = this.jwtAuthService.login({
      id: tokenInfo.userId!,
      login: tokenInfo.userName!,
      scopes: tokenInfo.scopes,
    });
    await this.authService.checkUser(code, tokenInfo.userId!, tokenInfo.userName);

    res.send({
      accessToken,
      refreshToken,
    });
  }

  @Post('token')
  async refresh(@Res() res: Express.Response, @Body() body: { refreshToken: string }) {
    if (!body.refreshToken) throw new BadRequestException('Refresh token not passed to body');
    try {
      const newTokens = await this.jwtAuthService.refresh(body.refreshToken);
      res.send(newTokens);
    } catch (error) {
      res.status(400).send('Something wrong with your authorization. Please try authorize again.');
    }
  }

  @Post('logout')
  @UseGuards(JwtAuthGuard)
  async logout(@Req() req: Express.Request) {
    await this.cacheManager.del(`nest:cache:auth/profile:${req.user.id}`);
    return true;
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
