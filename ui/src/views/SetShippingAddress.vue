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
                      Set Shipping Address
                    </DialogTitle>
                  </div>
                </div>
                <div class="overflow-y-auto flex-grow mt-4">
                  <form>
                    <div class="flex flex-col mb-4">
                      <label for="first-name">First Name</label>
                      <input
                        id="first-name"
                        v-model="form.firstName"
                        type="text"
                        placeholder="Enter your first name"
                        required
                      />
                    </div>
                    <div class="flex flex-col mb-4">
                      <label for="last-name">Last Name</label>
                      <input
                        id="last-name"
                        v-model="form.lastName"
                        type="text"
                        placeholder="Enter your last name"
                        required
                      />
                    </div>
                    <div class="flex flex-col mb-4">
                      <label for="address">Address</label>
                      <input
                        id="address"
                        v-model="form.address"
                        type="text"
                        placeholder="Enter your address"
                        required
                      />
                    </div>
                    <div class="flex flex-col mb-4">
                      <label for="address">Address2(Optional)</label>
                      <input
                        id="address"
                        v-model="form.address2"
                        type="text"
                        placeholder="Enter your address"
                      />
                    </div>
                    <!--div class="flex flex-col mb-4">
                      <label for="city">City</label>
                      <input
                        id="city"
                        v-model="form.city"
                        type="text"
                        placeholder="Enter your city"
                        required
                      />
                    </div>
                    <div class="flex flex-col mb-4">
                      <label for="state">State/Province</label>
                      <input
                        id="state"
                        v-model="form.state"
                        type="text"
                        placeholder="Enter your state or province"
                        required
                      />
                    </div-->
                    <div>
                      <label for="continent">Select Continent:</label>
                      <select
                        id="continent"
                        v-model="locationStore.selectedContinent"
                        @change="onContinentChange"
                      >
                        <option
                          v-for="continent in locationStore.continents"
                          :key="continent"
                          :value="continent"
                        >
                          {{ continent }}
                        </option>
                      </select>
                    </div>

                    <div v-if="locationStore.countries.length > 0">
                      <label for="country">Select Country:</label>
                      <select
                        id="country"
                        v-model="locationStore.selectedCountry"
                        @change="onCountryChange"
                      >
                        <option
                          v-for="country in locationStore.countries"
                          :key="country"
                          :value="country"
                        >
                          {{ country }}
                        </option>
                      </select>
                    </div>

                    <div v-if="locationStore.cities.length > 0">
                      <label for="city">Select City:</label>
                      <select
                        id="city"
                        v-model="locationStore.selectedCity"
                        @change="onCityChange"
                      >
                        <option
                          v-for="city in locationStore.cities"
                          :key="city"
                          :value="city"
                        >
                          {{ city }}
                        </option>
                      </select>
                    </div>
                    <div class="flex flex-col mb-4">
                      <label for="zip">Zip/Postal Code</label>
                      <input
                        id="zip"
                        v-model="form.zip"
                        type="text"
                        placeholder="Enter your zip or postal code"
                        required
                      />
                    </div>
                    <!--div class="flex flex-col mb-4">
                      <label for="country">Country</label>
                      <input
                        id="country"
                        v-model="form.country"
                        type="text"
                        placeholder="Enter your country"
                        required
                      />
                    </div-->
                    <div class="flex flex-col mb-4">
                      <label for="phone">Phone Number</label>
                      <input
                        id="phone"
                        v-model="form.phone"
                        type="tel"
                        placeholder="Enter your phone number"
                        required
                      />
                    </div>
                    <div class="flex flex-col mb-4">
                      <label for="email">Email</label>
                      <input
                        id="email"
                        v-model="form.email"
                        type="email"
                        placeholder="Enter your phone number"
                        required
                      />
                    </div>
                  </form>
                </div>
                <div class="bg-gray-50 px-4 py-3 flex flex-row-reverse sm:px-6">
                  <button
                    type="button"
                    class="inline-flex w-full justify-center rounded-md bg-blue px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-red-500 sm:ml-3 sm:w-auto"
                    @click="submitForm"
                  >
                    Finish
                  </button>
                  <button
                    type="button"
                    class="inline-flex w-full justify-center rounded-md bg-red-600 px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:mt-0 sm:w-auto"
                    @click="dialogStore.closeDialog"
                    ref="cancelButtonRef"
                    style="background-color: red"
                  >
                    Cancel
                  </button>
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
import { useSAStore } from "../store/dstore";
import { useLocationStore } from "../store/locationStore";
const locationStore = useLocationStore();
import {
  Dialog,
  DialogPanel,
  DialogTitle,
  TransitionChild,
  TransitionRoot,
} from "@headlessui/vue";
import { ArrowLeftIcon } from "@heroicons/vue/24/outline";
import { reactive } from "vue";
import { onMounted } from "vue";
import axios from "axios";
import { useNotification } from "@kyvg/vue3-notification";

const notification = useNotification();
const form = reactive({
  firstName: "",
  lastName: "",
  address: "",
  city: "",
  state: "",
  zip: "",
  country: "",
  phone: "",
  continent: "",
  address2: "",
  email: "",
});

onMounted(async () => {
  await locationStore.fetchContinents();
});
const submitForm = async () => {
  if (
    form.firstName &&
    form.lastName &&
    form.address &&
    form.city &&
    form.continent &&
    form.email &&
    form.zip &&
    form.country &&
    form.phone
  ) {
    try {
      const user = window.Telegram.WebApp.initDataUnsafe.user;
      await axios.post(
        `/api/v1/submitShippingAddress?user_id=${user.id}`,
        form
      );
      //alert(response.data.message);
      notification.notify({
        title: "Success",
        text: "Shipping address set successfully",
        type: "success",
      });
      dialogStore.closeDialog();
    } catch (error) {
      notification.notify({
        title: "Error",
        text: "error setting shipping address",
        type: "error",
      });
    }
  } else {
    notification.notify({
      title: "Incomplete information",
      text: "Please fill all the fileds",
      type: "warn",
    });
  }
};

const dialogStore = useSAStore();
const onContinentChange = () => {
  locationStore.setSelectedContinent();
  form.continent = locationStore.selectedContinent;
};

const onCountryChange = () => {
  locationStore.setSelectedCountry();
  form.country = locationStore.selectedCountry;
};

const onCityChange = () => {
  form.city = locationStore.selectedCity;
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
