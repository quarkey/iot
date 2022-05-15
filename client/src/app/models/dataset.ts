import { ChartDataset } from "chart.js";

export interface Dataset {
  id: number;
  sensor_id: number;
  title: string;
  description: string;
  reference: string;
  intervalsec: number;
  fields: string[];
  created_at: string;
  types: string[];
  showcharts: boolean[];
  sensor_title: string;
  telemetry?: string;
}

export interface Sensordata {
  id: number;
  data: string[];
  time: string;
  dataset_id?: number;
}

export interface Ng2Dataset {
  labels: string[];
  lineChartdataset: ChartDataset[];
}
