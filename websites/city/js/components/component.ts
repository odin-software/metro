import Store from "../store/store.js";

interface ComponentProps<T> {
  store?: Store<T>;
  element?: HTMLElement;
}

export class Component<T> {
  element?: HTMLElement;

  constructor(props: ComponentProps<T>) {
    let self = this;

    this.render = this.render || function () {};

    if (props.store instanceof Store) {
      props.store.events.subscribe("stateChange", () => self.render());
    }

    if (props.hasOwnProperty("element")) {
      this.element = props.element;
    }
  }

  render() {}
}
