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
      const title = this.element.querySelector(".confirm-dialog-body h1");
      const body = this.element.querySelector(".confirm-dialog-body p");
      const yesButton = this.element.querySelector(
        ".confirm-dialog-body #yesBtn"
      );
      const noButton = this.element.querySelector(
        ".confirm-dialog-body #yesBtn"
      );

      title.innerHTML = DialogStoreValue.state.title;
      body.innerHTML = DialogStoreValue.state.body;
      yesButton.addEventListener("click", () =>
        DialogStoreValue.state.yesBtn()
      );
      noButton.addEventListener("click", () => DialogStoreValue.state.noBtn());

      if (DialogStoreValue.state.open) {
        this.element.classList.add("open");
      } else {
        this.element.classList.remove("open");
      }
    }
  }
}
