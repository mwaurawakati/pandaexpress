// dialogStore.ts
import { defineStore } from "pinia";

interface DialogState {
  open: boolean;
}
export const useLanguageStore = defineStore("lang", {
  state: (): DialogState => ({
    open: false,
  }),
  actions: {
    openDialog() {
      this.open = true;
    },
    closeDialog() {
      this.open = false;
    },
  },
});

export const useSAStore = defineStore("d", {
  state: (): DialogState => ({
    open: false,
  }),
  actions: {
    openDialog() {
      this.open = true;
    },
    closeDialog() {
      this.open = false;
    },
  },
});
