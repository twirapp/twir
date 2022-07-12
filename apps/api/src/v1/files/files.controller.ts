import { Body, Controller, Delete, Param, Post, Req, UploadedFiles, UseGuards, UseInterceptors } from '@nestjs/common';
import { FilesInterceptor } from '@nestjs/platform-express';
import { Request } from 'express';

import { JwtAuthGuard } from '../../jwt/jwt.guard.js';
import { FilesService } from './files.service.js';

@Controller('v1/files')
export class FilesController {
  constructor(private readonly service: FilesService) { }

  @Post()
  @UseGuards(JwtAuthGuard)
  @UseInterceptors(FilesInterceptor('files'))
  uploadFile(@Req() req: Request, @UploadedFiles() files: Array<Express.Multer.File>) {
    return this.service.uploadFiles(files, req.user?.id);
  }

  @Delete(':id')
  @UseGuards(JwtAuthGuard)
  deleteFile(@Req() req: Request, @Param('id') id: string) {
    return this.service.deleteFile(id, req.user?.id);
  }
}
