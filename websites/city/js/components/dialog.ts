import { DialogStore } from "../typings/store.js";
import { Component } from "./component.js";
import DgStore from "../store/dialog.js";

export class Dialog extends Component<DialogStore> {
  constructor() {
    super({
      store: DgStore,
      element: document.querySelector("#dialog"),
    });
  }

  render() {
    console.log(DgStore);
  }
}
