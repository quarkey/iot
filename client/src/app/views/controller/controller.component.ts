import { Component, OnInit } from "@angular/core";
import { ActivatedRoute } from "@angular/router";
import { Controller } from "src/app/models/controllers";
import { ControllersService } from "src/app/services/controllers.service";

@Component({
  selector: "app-controller",
  templateUrl: "./controller.component.html",
  styleUrls: ["./controller.component.scss"],
})
export class ControllerComponent implements OnInit {
  citem: Controller;
  constructor(
    private activeRoute: ActivatedRoute,
    private ControllersService: ControllersService
  ) {}
  ngOnInit(): void {
    this.activeRoute.params.subscribe((key) => {
      this.ControllersService.LoadControllerByID(key.id).subscribe((res) => {
        this.citem = res;
      });
    });
  }
}
