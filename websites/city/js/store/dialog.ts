import Store from "./store.js";
import { DialogStore } from "../typings/store.js";

export const dialogStoreParams = {
  state: {
    open: false,
    title: "",
    body: "",
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
        open: payload.open,
        title: payload.title,
        body: payload.body,
      };
    },
  },
};

export default new Store<DialogStore>(dialogStoreParams);
