
import { HttpException, Injectable } from '@nestjs/common';
import { config } from '@tsuwari/config';

import { FeedBackPostDto } from './dto/post.dto.js';

@Injectable()
export class FeedbackService {
  async postFeedBack(data: FeedBackPostDto, userId: string) {
    if (!data.text) throw new HttpException(`Text is empty`, 400);

    const filesString = data.files?.map(f => `${config.MINIO_URL}/tsuwari/${userId}/${f.id}`).join('\n');

    await fetch('https://discord.com/api/webhooks/996446327263219762/mD4_zBAEDbwLuEsPlm9IremBt6393kV37TUKAN-Hq0CiCdjbh4ic2VZiEIuKDasLgaDI', {
      method: 'POST',
      body: JSON.stringify({
        'username': 'FeedBack from site',
        'content': `${userId}\n${data.text}\n${filesString ?? ''}`,
      }),
      headers: {
        'Content-type': 'application/json',
      },
    });

    return '';
  }
}
