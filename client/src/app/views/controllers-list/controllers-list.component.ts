import { Component, OnInit } from "@angular/core";
import { Controller } from "src/app/models/controllers";
import { ControllersService } from "src/app/services/controllers.service";
import { DialogsService } from "src/app/services/dialogs.service";

@Component({
  selector: "app-controllers",
  templateUrl: "./controllers-list.component.html",
  styleUrls: ["./controllers-list.component.scss"],
})
export class ControllersListComponent implements OnInit {
  citem: Controller[];
  loading: boolean = true;
  constructor(
    private ControllersService: ControllersService,
    private dialogService: DialogsService
  ) {}

  ngOnInit(): void {
    this.ControllersService.LoadControllersList().subscribe((res) => {
      if (res) {
        this.citem = res;
        this.loading = false;
      }
    });
  }
  newControllerDialog() {
    this.dialogService.openNewControllerDialog();
  }
}
