export interface Dataset {
  id: number;
  sensor_id: number;
  title: string;
  description: string;
  reference: string;
  intervalsec: number;
  fields: any;
  created_at: string;
  types: any;
}
