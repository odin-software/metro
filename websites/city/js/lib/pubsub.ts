export default class PubSub {
  events: Record<string, Function[]>;

  constructor() {
    this.events = {};
  }

  subscribe(event: string, callback) {
    let self = this;
    if (!self.events.hasOwnProperty(event)) {
      self.events[event] = [];
    }

    return self.events[event].push(callback);
  }

  publish(event: string, data = {}) {
    let self = this;
    if (!self.events.hasOwnProperty(event)) {
      return [];
    }

    return self.events[event].map((callback) => callback(data));
  }
}
