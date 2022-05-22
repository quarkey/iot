export interface Controller {
  id: number;
  sensor_id: number;
  title: string;
  description: string;
  switch: number;
  items: any;
  alert: boolean;
  active: boolean;
  created_at: string;
}