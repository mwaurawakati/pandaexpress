// src/stores/locationStore.ts
import { defineStore } from "pinia";
import { useNotification } from "@kyvg/vue3-notification";
import axios from "axios";

const notification = useNotification();

export const useLocationStore = defineStore("location", {
  state: () => ({
    continents: [],
    countries: [] as string[],
    cities: [] as string[],
    selectedContinent: "",
    selectedCountry: "",
    selectedCity: "",
  }),
  actions: {
    async fetchContinents() {
      try {
        const response = await axios.get("/api/v1/continents"); // Replace with your endpoint
        this.continents = response.data;
      } catch (error) {
        notification.notify({
          title: "Error",
          text: "Error fetching continents",
          type: "error",
        });
      }
    },
    async fetchCountries(continent: string) {
      try {
        const response = await axios.get(
          `/api/v1/countries?continent=${continent}`
        ); // Replace with your endpoint
        this.countries = response.data;
      } catch (error) {
        notification.notify({
          title: "Error",
          text: "Error fetching countries",
          type: "error",
        });
      }
    },
    async fetchCities(country: string) {
      try {
        const response = await axios.get(`/api/v1/cities?country=${country}`); // Replace with your endpoint
        this.cities = response.data;
      } catch (error) {
        notification.notify({
          title: "Error",
          text: "Error fetching cities",
          type: "error",
        });
      }
    },
    setSelectedContinent() {
      this.selectedCountry = "";
      this.selectedCity = "";
      this.countries = [];
      this.cities = [];
      this.fetchCountries(this.selectedContinent);
    },
    setSelectedCountry() {
      this.selectedCity = "";
      this.cities = [];
      this.fetchCities(this.selectedCountry);
    },
  },
});
