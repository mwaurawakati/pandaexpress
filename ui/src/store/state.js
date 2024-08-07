// state.js
import { ref } from "vue";

const open = ref(false);

function openDialog() {
  open.value = true;
}

function closeDialog() {
  open.value = false;
}

export function useDialog() {
  return {
    open,
    openDialog,
    closeDialog,
  };
}
