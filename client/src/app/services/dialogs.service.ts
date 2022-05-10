import { Injectable } from "@angular/core";
import {
  MatDialog,
  MatDialogRef,
  MAT_DIALOG_DATA,
} from "@angular/material/dialog";
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
}
