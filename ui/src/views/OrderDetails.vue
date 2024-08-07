<template>
  <div v-if="order">
    <h1>Order Details</h1>
    <div>
      <p><strong>Order ID:</strong> {{ order.ID }}</p>
      <p><strong>User ID:</strong> {{ order.UserID }}</p>
      <p><strong>Total Price:</strong> {{ order.TotalPrice }}</p>
      <p><strong>Shipping Fee:</strong> {{ order.ShippingFee }}</p>
      <p><strong>Refferee:</strong> {{ order.Refferee }}</p>
      <h2>Items</h2>
      <ul>
        <li v-for="item in order.Items" :key="item.ItemID">
          <img :src="item.Image" alt="Item Image" />
          <p><strong>Title:</strong> {{ item.ItemTitle }}</p>
          <p><strong>Price:</strong> {{ item.Price }}</p>
          <p><strong>Quantity:</strong> {{ item.Quantity }}</p>
        </li>
      </ul>
    </div>
  </div>
</template>
<script lang="ts">
import { defineComponent, onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import axios from "axios";

export default defineComponent({
  name: "OrderDetails",
  setup() {
    const route = useRoute();
    const order = ref<any>(null);

    onMounted(async () => {
      try {
        const orderId = route.params.id;
        const response = await axios.get(`/api/orders/${orderId}`);
        order.value = response.data;
      } catch (error) {
        console.error("Error fetching order details:", error);
      }
    });

    return { order };
  },
});
</script>
