import { TwitchAuthGuardOptions } from '@nestjs-hybrid-auth/twitch';
import { BadRequestException, Body, CacheTTL, Controller, ExecutionContext, Get, Post, Query, Req, Res, UseGuards, UseInterceptors } from '@nestjs/common';
import { Injectable } from '@nestjs/common';
import { AuthGuard, IAuthModuleOptions } from '@nestjs/passport';
import { config } from '@tsuwari/config';
import { exchangeCode, getTokenInfo } from '@twurple/auth';
import type { Request, Response } from 'express';
import merge from 'lodash.merge';

import { CustomCacheInterceptor } from '../helpers/customCacheInterceptor.js';
import { JwtAuthGuard } from '../jwt/jwt.guard.js';
import { JwtAuthService } from '../jwt/jwt.service.js';
import { AuthService } from './auth.service.js';

@Injectable()
class TwitchAuthGuard extends AuthGuard('twitch') {
  constructor(options?: TwitchAuthGuardOptions) {
    super(
      merge({}, options, {
        property: 'hybridAuthResult',
      }),
    );
  }

  getAuthenticateOptions(context: ExecutionContext): Promise<IAuthModuleOptions> | IAuthModuleOptions | undefined {
    const req = context.switchToHttp().getRequest() as Request;

    return req.query;
  }
}

function UseTwitchAuth(options?: TwitchAuthGuardOptions) {
  return UseGuards(new TwitchAuthGuard(options));
}

@Controller('auth')
export class AuthController {
  constructor(private readonly jwtAuthService: JwtAuthService, private readonly authService: AuthService) { }

  @UseTwitchAuth()
  @Get('')
  login() {
    return 'Login with twitch...';
  }

  @Get('token')
  async callback(@Res() res: Response, @Query() query: { code: string; state: string }) {
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
    });
    this.authService.checkUser(code, tokenInfo.userId!, tokenInfo.userName);

    res.send({
      accessToken,
      refreshToken,
    });
  }

  @Post('token')
  async refresh(@Body() body: { refreshToken: string }) {
    if (!body.refreshToken) throw new BadRequestException('Refresh token not passed to body');
    const newTokens = await this.jwtAuthService.refresh(body.refreshToken);

    return newTokens;
  }

  @CacheTTL(600)
  @UseInterceptors(CustomCacheInterceptor(ctx => {
    const req = ctx.switchToHttp().getRequest() as Request;
    return `nest:cache:auth/profile:${req.user.id}`;
  }))
  @UseGuards(JwtAuthGuard)
  @Get('profile')
  public showProfile(@Req() req: Request) {
    return this.authService.getProfile(req.user.id);
  }
}
