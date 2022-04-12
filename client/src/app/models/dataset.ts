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
  sensor_title: string;
}
