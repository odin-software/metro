interface EventMap {
  change: InputEvent;
}

interface EventEmitter {
  addEventListener<E extends keyof EventMap>(
    type: E,
    listener: (ev: EventMap[E]) => any
  ): void;
}

interface EventTarget {
  value: any;
}
