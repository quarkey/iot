import { Injectable } from "@angular/core";
import {
  MatDialog,
  MatDialogRef,
  MAT_DIALOG_DATA,
} from "@angular/material/dialog";
import { Observable, Subject } from "rxjs";
import { subscribeOn } from "rxjs/operators";
import { ConfirmationDialogComponent } from "../components/dialogs/confirmation-dialog/confirmation-dialog.component";
import { NewControllerComponent } from "../components/dialogs/new-controller/new-controller.component";
import { NewDatasetDialogComponent } from "../components/dialogs/new-dataset/new-dataset.component";
import { NewDeviceDialogComponent } from "../components/dialogs/new-device/new-device.component";

@Injectable({
  providedIn: "root",
})
export class DialogsService {
  constructor(public dialog: MatDialog) {}

  openNewDatasetDialog() {
    const dialogRef = this.dialog.open(NewDatasetDialogComponent, {
      width: "250px",
      data: {},
    });
    dialogRef.afterClosed().subscribe((result) => {
      console.log("The dialog was closed");
    });
  }
  openNewDeviceDialog() {
    const dialogRef = this.dialog.open(NewDeviceDialogComponent, {
      width: "250px",
      data: {},
    });
    dialogRef.afterClosed().subscribe((result) => {
      console.log("The dialog was closed");
    });
  }
  openNewControllerDialog() {
    const dialogRef = this.dialog.open(NewControllerComponent, {
      width: "250px",
      data: {},
    });
    dialogRef.afterClosed().subscribe((result) => {
      console.log("The dialog was closed");
    });
  }
  openConfirmationDialog(
    title: string,
    message: string,
    yesbutton?: string,
    nobutton?: string
  ): Observable<boolean> {
    const dialogRef = this.dialog.open(ConfirmationDialogComponent, {
      data: {
        title: title,
        message: message,
        yes: yesbutton || "CONFIRM",
        cancel: nobutton || "CANCEL",
      },
    });
    var out = new Subject<boolean>();
    dialogRef.afterClosed().subscribe((res) => {
      if (res == "yes") {
        out.next(true);
      } else {
        out.next(false);
      }
    });
    return out.asObservable();
  }
}
