import { TOUCH_BUFFER_MS } from "@angular/cdk/a11y/input-modality/input-modality-detector";
import { Component, Input, OnInit } from "@angular/core";
import { FormArray, FormBuilder, FormGroup, Validators } from "@angular/forms";
import { Dataset } from "src/app/models/dataset";
import { DatasetsService } from "src/app/services/datasets.service";

@Component({
  selector: "app-dataset-details",
  templateUrl: "./details.component.html",
  styleUrls: ["./details.component.scss"],
})
export class DetailsComponent implements OnInit {
  @Input() dataset: Dataset;
  form: FormGroup;

  constructor(
    private formBuilder: FormBuilder,
    private datasetService: DatasetsService
  ) {}

  ngOnInit(): void {
    this.form = this.formBuilder.group({
      title: [this.dataset.title, Validators.required],
      description: [this.dataset.description, Validators.required],
      intervalsec: [this.dataset.intervalsec, Validators.required],
      types: this.formBuilder.array([]),
      fields: this.formBuilder.array([]),
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
  }
  get types() {
    return this.form.get("types") as FormArray;
  }
  get fields() {
    return this.form.get("fields") as FormArray;
  }
  addTypeField() {
    // fields comes in pairs
    this.types.push(this.formBuilder.control(""));
    this.fields.push(this.formBuilder.control(""));
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
}
