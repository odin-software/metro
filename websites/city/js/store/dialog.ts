import Store from "./store.js";
import { DialogStore } from "../typings/store.js";

export const dialogStoreParams = {
  state: {
    open: false,
  },
  actions: {
    openDialog: (context: Store<DialogStore>, payload) => {
      context.commit("openDialog", payload);
    },
  },
  mutations: {
    openDialog: (state, payload) => {
      return {
        ...state,
        open: payload,
      };
    },
  },
};

export default new Store<DialogStore>(dialogStoreParams);
