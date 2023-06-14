import { Component, OnInit } from '@angular/core';
import { DatasetsService } from 'src/app/services/datasets.service';

@Component({
  selector: 'app-report-details',
  templateUrl: './report-details.component.html',
  styleUrls: ['./report-details.component.scss'],
})
export class ReportDetailsComponent implements OnInit {
  constructor(private datasetService: DatasetsService) {}
  report: any;
  ngOnInit(): void {
    this.datasetService.GetTemperatureReport().subscribe((res) => {
      if (res) {
        this.report = res;
      }
    });
  }
}
