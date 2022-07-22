type Opts = {
  interval: number,
}

type Cb = (id: string) => any | Promise<any>;

export class Queue<T> {
  #tasks: Map<string, {
    item: T,
    opts: Opts,
    data: {
      count: number
      lastTimeout: NodeJS.Timeout | null
    },
  }> = new Map();

  constructor(private readonly callback: Cb) { }

  addTimerToQueue(id: string, item: T, opts: Opts) {
    if (this.#tasks.has(id)) {
      throw new Error(`Task with id ${id} already exists.`);
    }

    this.#tasks.set(id, {
      item,
      opts,
      data: { count: 0, lastTimeout: null },
    });
    this.#process(id);
  }

  removeTimerFromQueue(id: string) {
    const task = this.#tasks.get(id);
    if (!task) return;
    if (task.data.lastTimeout) clearTimeout(task.data.lastTimeout);

    this.#tasks.delete(id);
    return;
  }

  async #process(id: string) {
    const task = this.#tasks.get(id);
    if (!task) return;
    task.data.count++;

    await this.callback(id);

    task.data.lastTimeout = setTimeout(() => this.#process(id), task?.opts.interval);
  }
}