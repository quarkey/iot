export interface Controller {
  id: number;
  sensor_id: number;
  category: string;
  title: string;
  description: string;
  switch: number;
  items: any;
  alert: boolean;
  active: boolean;
  created_at: string;
}
export interface thresholdswitch {
  on: boolean;
  item_description: string;
  datasource: string;
  threshold_limit: number;
  operation: string;
}
export interface timeswitch {
  on: boolean;
  item_description: string;
  time_on: string;
  time_off: string;
  duration: string;
  repeat: number;
}

export interface normalswitch {
  on: boolean;
  item_description: string;
}
