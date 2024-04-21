import PubSub from "../lib/pubsub.js";

export default class Store<T> {
  state: T;
  actions: Record<string, Function>;
  mutations: Record<string, Function>;
  events: PubSub;

  constructor(params) {
    let self = this;

    self.actions = {};
    self.mutations = {};
    self.state = params.state;

    self.events = new PubSub();

    if (params.hasOwnProperty("actions")) {
      self.actions = params.actions;
    }
    if (params.hasOwnProperty("mutations")) {
      self.mutations = params.mutations;
    }
  }

  setState(newState: T) {
    // TODO: Fix this so it can deep clone the object with
    // `structuredClone` and attach the functions to the state.
    this.state = newState;
    this.events.publish("stateChange", this.state);
  }

  dispatch(actionKey: string, payload): boolean {
    if (typeof this.actions[actionKey] !== "function") {
      console.error(`Action "${actionKey} doesn't exists.`);
      return false;
    }
    this.actions[actionKey](this, payload);
    return true;
  }

  commit(mutationKey: string, payload): boolean {
    if (typeof this.mutations[mutationKey] !== "function") {
      console.error(`Mutation "${mutationKey} doesn't exists.`);
      return false;
    }
    let newState = this.mutations[mutationKey](this.state, payload);
    this.setState(newState);
    return true;
  }
}
