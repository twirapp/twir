import { randomUUID } from 'node:crypto';

import { HttpException, Injectable } from '@nestjs/common';
import { PrismaService } from '@tsuwari/prisma';
import S3 from 'nestjs-s3';

@Injectable()
export class FilesService {
  constructor(
    @S3.InjectS3() private readonly s3: S3.S3,
    private readonly prisma: PrismaService,
  ) { }

  async uploadFiles(files: Array<Express.Multer.File>, channelId?: string) {
    const result = await Promise.all(files.map((file) => {
      const id = randomUUID();

      return Promise.all([
        this.prisma.file.create({
          data: {
            id,
            name: file.originalname,
            userId: channelId,
            size: file.size,
            type: file.mimetype,
          },
        }),
        this.s3.upload({
          Bucket: 'tsuwari',
          Key: `${channelId}/${id}`,
          ACL: 'public-read',
          Body: file.buffer,
          ContentType: `${file.mimetype}; charset=utf-8`,
        }).promise(),
      ]);
    }));

    return result.map(r => r[0]);
  }

  async deleteFile(id: string, userId: string) {
    const [file, dashboardAccess] = await Promise.all([
      this.prisma.file.findFirst({ where: { id } }),
      this.prisma.dashboardAccess.findMany({ where: { userId } }),
    ]);

    if (!file) throw new HttpException(`File with id ${id} not found.`, 404);
    if (!file.userId) throw new HttpException(`You cannot delete this file.`, 401);
    if (!dashboardAccess.find(a => a.channelId === file.userId) && file.userId !== userId) {
      throw new HttpException(`You cannot delete this file.`, 401);
    }


    const [deletedFile] = await Promise.all([
      this.prisma.file.delete({ where: { id } }),
      this.s3.deleteObject({ Bucket: 'tsuwari', Key: `${file.id}/${id}` }),
    ]);
    return deletedFile;
  }
}
