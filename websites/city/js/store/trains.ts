import Store from "./store";
import { TrainStore } from "../typings/store";

export const trainStoreParams = {
  state: {
    trains: [],
  },
  actions: {
    updateAllTrains: (context: Store<TrainStore>, payload) => {
      context.commit("updateAllTrains", payload);
    },
  },
  mutations: {
    updateAllTrains: (state, payload) => {
      return payload;
    },
  },
};
