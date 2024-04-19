import { DialogStore } from "../typings/store.js";
import { Component } from "./component.js";
import DialogStoreValue from "../store/dialog.js";

export class Dialog extends Component<DialogStore> {
  constructor() {
    super({
      store: DialogStoreValue,
      element: document.querySelector("#dialog"),
    });

    this.render();
  }

  render() {
    if (this.element) {
      const title = this.element.querySelector(".jw-modal-body h1");
      title.innerHTML = DialogStoreValue.state.title;
      if (DialogStoreValue.state.open) {
        this.element.classList.add("open");
      } else {
        this.element.classList.remove("open");
      }
    }
    console.log(DialogStoreValue);
  }
}
