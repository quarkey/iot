import { Component, Input, OnInit } from '@angular/core';
import { FormArray, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { Dataset, Sensordata } from 'src/app/models/dataset';
import { DatasetsService } from 'src/app/services/datasets.service';
import { DialogsService } from 'src/app/services/dialogs.service';
import { GeneralService } from 'src/app/services/general.service';

@Component({
  selector: 'app-dataset-details',
  templateUrl: './dataset-details.component.html',
  styleUrls: ['./dataset-details.component.scss'],
})
export class DatasetDetailsComponent implements OnInit {
  @Input() dataset: Dataset;
  form: FormGroup;
  loadingdownloadFile = false;
  liveSensordata: Sensordata;
  constructor(
    private formBuilder: FormBuilder,
    private datasetService: DatasetsService,
    private generalService: GeneralService,
    private dialogService: DialogsService,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.form = this.formBuilder.group({
      title: [this.dataset.title, Validators.required],
      description: [this.dataset.description, Validators.required],
      intervalsec: [this.dataset.intervalsec, Validators.required],
      types: this.formBuilder.array([]),
      fields: this.formBuilder.array([]),
      showcharts: this.formBuilder.array([]),
    });
    this.populateFormArray();
  }
  populateFormArray() {
    this.dataset.types.forEach((x) => {
      this.types.push(this.formBuilder.control(x));
    });
    this.dataset.fields.forEach((x) => {
      this.fields.push(this.formBuilder.control(x));
    });
    if (this.dataset.showcharts === null) {
      for (let i = 0; i < this.dataset.fields.length; i++) {
        this.showcharts.push(this.formBuilder.control(false));
      }
    } else {
      this.dataset.showcharts.forEach((x) => {
        this.showcharts.push(this.formBuilder.control(x));
      });
    }
  }
  get types() {
    return this.form.get('types') as FormArray;
  }
  get fields() {
    return this.form.get('fields') as FormArray;
  }
  get showcharts() {
    return this.form.get('showcharts') as FormArray;
  }
  addTypeField() {
    // fields comes in pairs
    this.types.push(this.formBuilder.control(''));
    this.fields.push(this.formBuilder.control(''));
    this.showcharts.push(this.formBuilder.control(false));
  }
  updateDataset() {
    var obj = {
      ...this.form.value,
      reference: this.dataset.reference,
      id: this.dataset.id,
    };

    this.datasetService.updateDataset(obj).subscribe((res) => {
      if (res) {
        this.form.markAsPristine();
      }
    });
  }
  downloadCSV() {
    this.loadingdownloadFile = true;
    this.datasetService.LoadCSVDatasetByReference(this.dataset.reference).subscribe((res) => {
      if (res) {
        this.loadingdownloadFile = false;
        const date = Date.now();
        const filename = `export_dataset_id_${this.dataset.id}_${date}.csv`;
        this.generalService.DownloadFile(res, filename);
      }
    });
  }
  deleteDataset() {
    this.dialogService
      .openConfirmationDialog(
        'Delete dataset, INCLUDING DATA?',
        `Are you sure you want to permanently delete dataset, INCLUDING data points?`,
        'CONFIRM',
        'CANCEL',
        true
      )
      .subscribe((confirm) => {
        if (confirm) {
          this.datasetService.DeleteDatasetByID(this.dataset.id, this.dataset.title).subscribe((res) => {
            if (res) {
              this.router.navigate([`/datasets`]);
            }
          });
        }
      });
  }
}
