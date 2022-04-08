import { Component, Input, OnInit } from "@angular/core";
import { FormBuilder, FormGroup, Validators } from "@angular/forms";
import { Dataset } from "src/app/models/dataset";

@Component({
  selector: "app-dataset-details",
  templateUrl: "./details.component.html",
  styleUrls: ["./details.component.scss"],
})
export class DetailsComponent implements OnInit {
  @Input() dataset: Dataset;
  form: FormGroup;

  constructor(private formBuilder: FormBuilder) {}

  ngOnInit(): void {}
}
