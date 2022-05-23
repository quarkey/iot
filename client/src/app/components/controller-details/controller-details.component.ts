import { Component, Input, OnInit } from "@angular/core";
import { FormArray, FormBuilder, FormGroup, Validators } from "@angular/forms";
import { Controller } from "src/app/models/controllers";
import { Sensordata } from "src/app/models/dataset";
import { ControllersService } from "src/app/services/controllers.service";
import { environment } from "src/environments/environment";

@Component({
  selector: "app-controller-details",
  templateUrl: "./controller-details.component.html",
  styleUrls: ["./controller-details.component.scss"],
})
export class ControllerDetailsComponent implements OnInit {
  @Input() citem: Controller;
  form: FormGroup;
  categories: string[] = ["switch", "thresholdswitch", "timeswitch"];
  socket: any;
  liveTelemetry: any;
  loading: boolean = false;
  datasource: any;
  constructor(
    private formBuilder: FormBuilder,
    private controllerService: ControllersService
  ) {}

  ngOnInit(): void {
    const thresholdForm = this.formBuilder.group({
      on: [this.citem.items[0].on],
      item_description: [
        this.citem.items[0].item_description,
        Validators.required,
      ],
      operation: [this.citem.items[0].operation, Validators.required],
      datasource: [this.citem.items[0].datasource, Validators.required],
      threshold_limit: [
        this.citem.items[0].threshold_limit,
        Validators.required,
      ],
    });

    this.form = this.formBuilder.group({
      category: [this.citem.category, Validators.required],
      title: [this.citem.title, Validators.required],
      description: [this.citem.description, Validators.required],
      items: thresholdForm,
      active: [this.citem.active],
    });
    this.datasource = this.citem.items[0].datasource;
    this.runLive(this.datasource);
    this.form.get("items").valueChanges.subscribe((x) => {
      this.runLive(this.datasource);
    });
  }
  updateController() {
    var obj = {
      ...this.form.value,
      id: this.citem.id,
    };
    var temp = [obj.items];
    obj.items = temp;

    this.controllerService.UpdateControllerByID(obj).subscribe((res) => {
      if (res) {
        this.form.markAsPristine();
      }
    });
  }
  runLive(datasource: string) {
    this.loading = true;
    this.liveTelemetry = "loading..";
    const socket = new WebSocket(`${environment.wsUrl}/api/live`);
    socket.onopen = function (e) {
      console.log("WebSocket Opened");
    };
    this.socket = socket;
    var self = this;
    const regexpSize = /d([0-9]+)c([0-9]+)/;
    const match = datasource.match(regexpSize);

    socket.onmessage = function (e) {
      const data = JSON.parse(e.data) as Sensordata;
      // only showing current dataset
      // console.log(data.dataset_id);
      // console.log(match[1]);
      if (data.dataset_id === parseInt(match[1])) {
        console.log(data.data[match[2]]);
        self.liveTelemetry = data.data[match[2]] | 0;
        self.loading = false;
      }
    };
  }
}
