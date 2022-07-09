import { HttpException, Injectable } from '@nestjs/common';

@Injectable()
export class VersionService {
  async getCommitSha() {
    const request = await fetch('https://api.github.com/repos/satont/tsuwari/commits');
    if (!request.ok) {
      console.log(request);
      throw new HttpException('Request not ok', 500);
    }

    const data = await request.json();

    return data[0].sha;
  }
}
