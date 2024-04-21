import { DialogStore } from "../typings/store.js";
import { Component } from "./component.js";
import DialogStoreValue from "../store/dialog.js";

export class Dialog extends Component<DialogStore> {
  yesF: () => void;
  noF: () => void;

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
        ".confirm-dialog-body #noBtn"
      );

      title.innerHTML = DialogStoreValue.state.title;
      body.innerHTML = DialogStoreValue.state.body;
      yesButton.removeEventListener("click", this.yesF);
      noButton.removeEventListener("click", this.noF);
      yesButton.addEventListener("click", DialogStoreValue.state.yesBtn);
      noButton.addEventListener("click", DialogStoreValue.state.noBtn);

      this.yesF = DialogStoreValue.state.yesBtn;
      this.noF = DialogStoreValue.state.yesBtn;

      if (DialogStoreValue.state.open) {
        this.element.classList.add("open");
      } else {
        this.element.classList.remove("open");
      }
    }
  }
}
