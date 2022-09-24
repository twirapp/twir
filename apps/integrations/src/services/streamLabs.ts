import * as IO from 'socket.io-client';

type Socket = typeof IO.Socket;

export class StreamLabs {
  #conn: Socket;

  constructor(token: string) {
    this.#conn = IO.connect(`https://sockets.streamlabs.com?token=${token}`, {
      transports: ['websocket'],
    });

    this.#conn.on('event', (eventData: Event) => {
      if (eventData.type === 'donation') {
        eventData.message.forEach((m) => this.#handler(m));
      }
    });
  }

  async #handler(data: Message) {
    return console.log(data);
  }
}

export interface Event {
  type: 'donation';
  event_id: string;
  message: Message[];
}

export interface Message {
  id: number;
  name: string;
  amount: string;
  formatted_amount: string;
  formattedAmount: string;
  message: string;
  currency: string;
  emotes?: any;
  iconClassName: string;
  to: {
    name: string;
  };
  from: string;
  from_user_id: null | string;
  _id: string;
}
