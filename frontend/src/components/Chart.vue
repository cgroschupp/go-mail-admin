<script setup>
import { Doughnut } from "vue-chartjs";
import Client from "../service/Client";
import { Chart as ChartJS, Title, ArcElement, Legend, Tooltip } from 'chart.js'

const props = defineProps({
  chartData: {
    type: Object,
    required: true
  },
  options: {
    type: Object,
    default: () => ({})
  }
})


ChartJS.register(Title, ArcElement, Legend, Tooltip)

</script>

<template>
  <Doughnut :data="chartData" :options="options" />
</template>