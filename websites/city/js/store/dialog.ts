import Store from "./store.js";
import { DialogStore } from "../typings/store.js";

export const dialogStoreParams = {
  state: {
    open: false,
    title: "",
    body: "",
    yesBtn: () => {},
    noBtn: () => {},
  },
  actions: {
    closeDialog: (context: Store<DialogStore>, payload) => {
      context.commit("closeDialog", payload);
    },
    openDialog: (context: Store<DialogStore>, payload) => {
      context.commit("openDialog", payload);
    },
  },
  mutations: {
    closeDialog: (state: DialogStore, _payload) => {
      return {
        ...state,
        open: false,
      };
    },
    openDialog: (state: DialogStore, payload) => {
      return {
        ...state,
        open: payload.open,
        title: payload.title,
        body: payload.body,
        yesBtn: payload.yesBtn,
        noBtn: payload.noBtn,
      };
    },
  },
};

export default new Store<DialogStore>(dialogStoreParams);
