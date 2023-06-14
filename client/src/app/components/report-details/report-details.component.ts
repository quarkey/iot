import { Component, Input, OnInit } from '@angular/core';
import { Dataset } from 'src/app/models/dataset';
import { DatasetsService } from 'src/app/services/datasets.service';

@Component({
  selector: 'app-report-details',
  templateUrl: './report-details.component.html',
  styleUrls: ['./report-details.component.scss'],
})
export class ReportDetailsComponent implements OnInit {
  constructor(private datasetService: DatasetsService) {}
  @Input() dataset: Dataset;
  report: any;
  ngOnInit(): void {
    this.datasetService.GetTemperatureReport(this.dataset.reference, 'd1c0').subscribe((res) => {
      if (res) {
        this.report = res;
      }
    });
  }
}
