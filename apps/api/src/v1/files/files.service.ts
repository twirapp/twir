import { randomUUID } from 'node:crypto';

import { HttpException, Injectable } from '@nestjs/common';
import { DashboardAccess } from '@tsuwari/typeorm/entities/DashboardAccess';
import { UserFile } from '@tsuwari/typeorm/entities/UserFile';
import S3 from 'nestjs-s3';

import { typeorm } from '../../index.js';

@Injectable()
export class FilesService {
  constructor(@S3.InjectS3() private readonly s3: S3.S3) {}

  async uploadFiles(files: Array<Express.Multer.File>, channelId?: string) {
    const result = await Promise.all(
      files.map((file) => {
        const id = randomUUID();

        return Promise.all([
          typeorm.getRepository(UserFile).save({
            id,
            name: file.originalname,
            userId: channelId,
            size: file.size,
            type: file.mimetype,
          }),
          this.s3
            .upload({
              Bucket: 'tsuwari',
              Key: `${channelId}/${id}`,
              ACL: 'public-read',
              Body: file.buffer,
              ContentType: `${file.mimetype}; charset=utf-8`,
            })
            .promise(),
        ]);
      }),
    );

    return result.map((r) => r[0]);
  }

  async deleteFile(id: string, userId: string) {
    const [file, dashboardAccess] = await Promise.all([
      typeorm.getRepository(UserFile).findOneBy({ id }),
      typeorm.getRepository(DashboardAccess).findBy({ userId }),
    ]);

    if (!file) throw new HttpException(`File with id ${id} not found.`, 404);
    if (!file.userId) throw new HttpException(`You cannot delete this file.`, 401);
    if (!dashboardAccess.find((a) => a.channelId === file.userId) && file.userId !== userId) {
      throw new HttpException(`You cannot delete this file.`, 401);
    }

    const [deletedFile] = await Promise.all([
      typeorm.getRepository(UserFile).delete({ id }),
      this.s3.deleteObject({ Bucket: 'tsuwari', Key: `${file.id}/${id}` }),
    ]);
    return deletedFile;
  }
}
