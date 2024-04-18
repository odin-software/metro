import PubSub from "../lib/pubsub.js";

export default class Store<T> {
  actions: Record<string, Function>;
  mutations: Record<string, Function>;
  state: T;
  status: string;
  events: PubSub;

  constructor(params) {
    let self = this;

    self.actions = {};
    self.mutations = {};
    self.state = params.state;
    self.status = "resting";

    self.events = new PubSub();

    if (params.hasOwnProperty("actions")) {
      self.actions = params.actions;
    }
    if (params.hasOwnProperty("mutations")) {
      self.mutations = params.mutations;
    }

    self.state = new Proxy(params.state || {}, {
      set(state, key, value, _) {
        state[key] = value;
        // console.log(`stateChange: ${String(key)}: ${value}`);

        self.events.publish("stateChange", self.state);

        if (self.status !== "mutation") {
          console.warn(
            `You should use a mutation to change set ${String(key)}`
          );
        }

        self.status = "resting";

        return true;
      },
    });
  }

  dispatch(actionKey: string, payload): boolean {
    let self = this;

    if (typeof self.actions[actionKey] !== "function") {
      console.error(`Action "${actionKey} doesn't exists.`);
      return false;
    }

    // console.groupCollapsed(`ACTION: ${actionKey}`);

    self.status = "action";
    self.actions[actionKey](self, payload);

    // console.groupEnd();

    return true;
  }

  commit(mutationKey: string, payload): boolean {
    let self = this;

    if (typeof self.mutations[mutationKey] !== "function") {
      console.error(`Mutation "${mutationKey} doesn't exists.`);
      return false;
    }

    self.status = "mutation";
    let newState = self.mutations[mutationKey](self.state, payload);
    self.state = Object.assign(self.state, newState);

    return true;
  }
}
