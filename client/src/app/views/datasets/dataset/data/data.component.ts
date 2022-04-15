import { Component, Input, OnInit } from "@angular/core";
import { Dataset, Sensordata } from "src/app/models/dataset";
import { DatasetsService } from "src/app/services/datasets.service";
import { multi } from "./data";

@Component({
  selector: "app-dataset-data",
  templateUrl: "./data.component.html",
  styleUrls: ["./data.component.scss"],
})
export class DataComponent implements OnInit {
  @Input() dataset: Dataset;
  data: any;
  loading: boolean = true;
  constructor(private datasetService: DatasetsService) {
    Object.assign(this, { multi });
  }

  multi: any[];
  view: any[] = [700, 300];

  // options
  legend: boolean = true;
  showLabels: boolean = true;
  animations: boolean = false;
  xAxis: boolean = true;
  yAxis: boolean = true;
  showYAxisLabel: boolean = true;
  showXAxisLabel: boolean = true;
  xAxisLabel: string = "Value";
  yAxisLabel: string = "Time";
  timeline: boolean = false;

  colorScheme = {
    domain: ["#5AA454", "#E44D25", "#CFC0BB", "#7aa3e5", "#a8385d", "#aae3f5"],
  };

  onSelect(event) {
    console.log(event);
  }
  ngOnInit(): void {
    this.datasetService
      .LoadAreaChartDatasetByReference(this.dataset.reference)
      .subscribe((res) => {
        this.data = res;
        if (res) {
          this.multi = res;
          this.loading = false;
        }
      });
  }
}
