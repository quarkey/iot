import { Component, Input, OnInit } from '@angular/core';
import { DatasetsService } from 'src/app/services/datasets.service';

@Component({
  selector: 'app-dataset-data',
  templateUrl: './data.component.html',
  styleUrls: ['./data.component.scss'],
})
export class DataComponent implements OnInit {
  @Input() reference: string;
  data: any;
  constructor(private datasetService: DatasetsService) {}

  ngOnInit(): void {
    this.datasetService
      .LoadDatasetByReference(this.reference)
      .subscribe((res) => {
        this.data = res;
      });
  }
}
