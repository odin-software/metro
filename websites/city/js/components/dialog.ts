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
      const title = this.element.querySelector(".dialog-body h1");
      const body = this.element.querySelector(".dialog-body p");
      const input = this.element.querySelector(
        ".dialog-body #dialogInput"
      ) as HTMLInputElement;
      const yesButton = this.element.querySelector(".dialog-body #yesBtn");
      const noButton = this.element.querySelector(".dialog-body #noBtn");

      title.innerHTML = DialogStoreValue.state.title;
      body.innerHTML = DialogStoreValue.state.body;
      yesButton.removeEventListener("click", this.yesF);
      noButton.removeEventListener("click", this.noF);
      yesButton.addEventListener("click", DialogStoreValue.state.yesBtn);
      noButton.addEventListener("click", DialogStoreValue.state.noBtn);

      this.yesF = DialogStoreValue.state.yesBtn;
      this.noF = DialogStoreValue.state.noBtn;

      if (DialogStoreValue.state.input.length > 0) {
        input.defaultValue = DialogStoreValue.state.input;
        input.addEventListener("change", (ev) => {
          //@ts-ignore
          DialogStoreValue.unreactive.input = ev.target.value;
        });
        input.classList.add("open");
      } else {
        input.classList.remove("open");
      }

      if (DialogStoreValue.state.open) {
        this.element.classList.add("open");
      } else {
        this.element.classList.remove("open");
      }
    }
  }
}
