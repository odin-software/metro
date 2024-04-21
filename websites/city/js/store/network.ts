import Store from "./store.js";
import { NetworkStore } from "../typings/store.js";
import { Network } from "../models/network.js";

export const networkStoreParams = {
  state: {
    network: new Network(),
  },
  actions: {},
  mutations: {},
};

export default new Store<NetworkStore>(networkStoreParams);
