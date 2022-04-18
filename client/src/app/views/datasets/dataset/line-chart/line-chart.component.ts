import { Component, Input, OnInit } from "@angular/core";
import { ChartDataset, ChartOptions, Color } from "chart.js";
import { Dataset, Ng2Dataset, Sensordata } from "src/app/models/dataset";
import { DatasetsService } from "src/app/services/datasets.service";

@Component({
  selector: "app-line-chart",
  templateUrl: "./line-chart.component.html",
  styleUrls: ["./line-chart.component.scss"],
})
export class LineChartComponent implements OnInit {
  @Input() dataset: Dataset;
  data: Ng2Dataset;
  constructor(private datasetService: DatasetsService) {}
  lineChartData: ChartDataset[];
  lineChartLabels: string[];

  lineChartOptions = {
    responsive: true,
    plugins: {
      legend: {
        position: "right",
      },
      title: {
        display: true,
        text: "Chart.js Line Chart",
      },
    },
  };
  lineChartLegend = true;
  lineChartPlugins = [];
  lineChartType = "line";

  ngOnInit(): void {
    this.datasetService
      .LoadLineChartDatasetByReference(this.dataset.reference)
      .subscribe((res) => {
        if (res) {
          this.lineChartLabels = res.labels;
          this.lineChartData = res.lineChartdataset;
        }
      });
  }
}
