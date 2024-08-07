<template>
  <TransitionRoot :show="dialogStore.open">
    <Dialog class="relative z-10" @close="dialogStore.closeDialog">
      <TransitionChild
        as="template"
        enter="ease-out duration-300"
        enter-from="opacity-0"
        enter-to="opacity-100"
        leave="ease-in duration-200"
        leave-from="opacity-100"
        leave-to="opacity-0"
      >
        <div
          class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"
        />
      </TransitionChild>

      <div class="fixed inset-0 z-10 w-screen overflow-y-auto">
        <div
          class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0"
        >
          <TransitionChild
            as="template"
            enter="ease-out duration-300"
            enter-from="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
            enter-to="opacity-100 translate-y-0 sm:scale-100"
            leave="ease-in duration-200"
            leave-from="opacity-100 translate-y-0 sm:scale-100"
            leave-to="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
          >
            <DialogPanel
              class="relative transform overflow-hidden rounded-lg bg-white text-left shadow-xl transition-all"
              style="margin: 5px; height: 90vh; width: 100vw"
            >
              <div
                class="bg-white px-4 pb-4 pt-5 sm:p-6 sm:pb-4 h-full flex flex-col"
              >
                <div class="flex sm:items-start">
                  <div
                    class="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-full bg-red-100 sm:mx-0 sm:h-10 sm:w-10"
                    @click="dialogStore.closeDialog"
                  >
                    <ArrowLeftIcon
                      class="h-8 w-8 text-blue-600"
                      aria-hidden="true"
                    />
                  </div>
                  <div
                    class="mx-auto mt-3 text-center sm:ml-4 sm:mt-0 sm:text-left"
                  >
                    <DialogTitle
                      as="h1"
                      class="text-2xl font-semibold leading-6 text-gray-900"
                    >
                      Set Preferred Language
                    </DialogTitle>
                  </div>
                </div>
                <div class="overflow-y-auto flex-grow mt-4">
                  <form @submit.prevent="submitForm">
                    <div>
                      <label for="language">Select Language:</label>
                      <select id="language" v-model="form.code" required>
                        <option
                          v-for="language in languages"
                          :key="language.Code"
                          :value="language.Code"
                        >
                          {{ language.Name }}
                        </option>
                      </select>
                    </div>
                    <div
                      class="bg-gray-50 px-4 py-3 flex flex-row-reverse sm:px-6"
                    >
                      <button
                        type="submit"
                        class="inline-flex w-full justify-center rounded-md bg-blue-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 sm:ml-3 sm:w-auto"
                      >
                        Finish
                      </button>
                      <button
                        type="button"
                        class="inline-flex w-full justify-center rounded-md bg-red-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-red-500 sm:mt-0 sm:w-auto"
                        @click="dialogStore.closeDialog"
                      >
                        Cancel
                      </button>
                    </div>
                  </form>
                </div>
              </div>
            </DialogPanel>
          </TransitionChild>
        </div>
      </div>
    </Dialog>
  </TransitionRoot>
</template>

<script setup>
import { reactive } from "vue";
import { useNotification } from "@kyvg/vue3-notification";
import axios from "axios";
import {
  Dialog,
  DialogPanel,
  DialogTitle,
  TransitionChild,
  TransitionRoot,
} from "@headlessui/vue";
import { ArrowLeftIcon } from "@heroicons/vue/24/outline";
import { useLanguageStore } from "../store/dstore"; // Assuming you have a store for managing dialogs

const dialogStore = useLanguageStore();
const notification = useNotification();

const languages = [
  { Name: "Afrikaans", Code: "af", EnglishName: "Afrikaans" },
  { Name: "አማርኛ", Code: "am", EnglishName: "Amharic" },
  { Name: "Български", Code: "bg", EnglishName: "Bulgarian" },
  { Name: "Català", Code: "ca", EnglishName: "Catalan" },
  {
    Name: "Chinese (Literary)",
    Code: "lzh",
    EnglishName: "Chinese (Literary)",
  },
  {
    Name: "Chinese Simplified",
    Code: "zh-Hans",
    EnglishName: "Chinese Simplified",
  },
  {
    Name: "Chinese Traditional",
    Code: "zh-Hant",
    EnglishName: "Chinese Traditional",
  },
  { Name: "Hrvatski", Code: "hr", EnglishName: "Croatian" },
  { Name: "Čeština", Code: "cs", EnglishName: "Czech" },
  { Name: "Dansk", Code: "da", EnglishName: "Danish" },
  { Name: "Nederlands", Code: "nl", EnglishName: "Dutch" },
  { Name: "English", Code: "en", EnglishName: "English" },
  { Name: "Eesti", Code: "et", EnglishName: "Estonian" },
  { Name: "Filipino", Code: "fil", EnglishName: "Filipino" },
  { Name: "Suomi", Code: "fi", EnglishName: "Finnish" },
  { Name: "Français (Canada)", Code: "fr-ca", EnglishName: "French (Canada)" },
  { Name: "Français (France)", Code: "fr", EnglishName: "French (France)" },
  { Name: "Deutsch", Code: "de", EnglishName: "German" },
  { Name: "Ελληνικά", Code: "el", EnglishName: "Greek" },
  { Name: "עברית", Code: "he", EnglishName: "Hebrew" },
  { Name: "हिन्दी", Code: "hi", EnglishName: "Hindi" },
  { Name: "Magyar", Code: "hu", EnglishName: "Hungarian" },
  { Name: "Íslenska", Code: "is", EnglishName: "Icelandic" },
  { Name: "Bahasa Indonesia", Code: "id", EnglishName: "Indonesian" },
  { Name: "Italiano", Code: "it", EnglishName: "Italian" },
  { Name: "日本語", Code: "ja", EnglishName: "Japanese" },
  { Name: "한국어", Code: "ko", EnglishName: "Korean" },
  { Name: "Latviešu", Code: "lv", EnglishName: "Latvian" },
  { Name: "Lietuvių", Code: "lt", EnglishName: "Lithuanian" },
  { Name: "Bahasa Melayu", Code: "ms", EnglishName: "Malay" },
  { Name: "Norsk", Code: "nb", EnglishName: "Norwegian" },
  { Name: "Polski", Code: "pl", EnglishName: "Polish" },
  {
    Name: "Português (Brasil)",
    Code: "pt",
    EnglishName: "Portuguese (Brazil)",
  },
  {
    Name: "Português (Portugal)",
    Code: "pt-pt",
    EnglishName: "Portuguese (Portugal)",
  },
  { Name: "Română", Code: "ro", EnglishName: "Romanian" },
  { Name: "Русский", Code: "ru", EnglishName: "Russian" },
  { Name: "Српски", Code: "sr-Latn", EnglishName: "Serbian" },
  { Name: "Slovenčina", Code: "sk", EnglishName: "Slovak" },
  { Name: "Slovenščina", Code: "sl", EnglishName: "Slovenian" },
  { Name: "Español (España)", Code: "es", EnglishName: "Spanish (Spain)" },
  { Name: "Kiswahili", Code: "sw", EnglishName: "Swahili" },
  { Name: "Svenska", Code: "sv", EnglishName: "Swedish" },
  { Name: "ไทย", Code: "th", EnglishName: "Thai" },
  { Name: "Türkçe", Code: "tr", EnglishName: "Turkish" },
  { Name: "Українська", Code: "uk", EnglishName: "Ukrainian" },
  { Name: "Tiếng Việt", Code: "vi", EnglishName: "Vietnamese" },
  { Name: "IsiZulu", Code: "zu", EnglishName: "Zulu" },
];

const form = reactive({
  code: "",
});

const submitForm = async () => {
  if (form.code) {
    try {
      const user = window.Telegram.WebApp.initDataUnsafe.user;
      await axios.post(`/api/v1/language?user_id=${user.id}`, form);
      notification.notify({
        title: "Success",
        text: "Preferred language set successfully",
        type: "success",
      });
      dialogStore.closeDialog();
    } catch (error) {
      notification.notify({
        title: "Error",
        text: "Error setting preferred language",
        type: "error",
      });
    }
  } else {
    notification.notify({
      title: "Incomplete information",
      text: "Please select a language",
      type: "warn",
    });
  }
};
</script>

<style scoped>
.dialog-wrapper {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}

.form-container {
  max-height: 60vh;
  overflow-y: auto;
  padding: 1rem;
}

form > div {
  margin-bottom: 1rem;
}

label {
  font-weight: bold;
  margin-bottom: 0.5rem;
  display: block;
}

input {
  padding: 0.5rem;
  border: 1px solid #ccc;
  border-radius: 4px;
  width: 100%;
}

button {
  padding: 0.75rem;
  border: none;
  border-radius: 4px;
  background-color: #007bff;
  color: white;
  cursor: pointer;
}

button:hover {
  background-color: #0056b3;
}
</style>
