import {
  HttpEvent,
  HttpInterceptor,
  HttpHandler,
  HttpRequest,
  HttpErrorResponse,
} from "@angular/common/http";
import { Injectable } from "@angular/core";
import { MatLegacyDialog as MatDialog } from "@angular/material/legacy-dialog";
import { Observable, throwError } from "rxjs";
import { catchError } from "rxjs/operators";
import { ErrorHandlingDialogComponent } from "../components/dialogs/error-handling-dialog/error-handling-dialog.component";
@Injectable()
export class HttpErrorInterceptor implements HttpInterceptor {
  constructor(private dialog: MatDialog) {}
  intercept(
    request: HttpRequest<any>,
    next: HttpHandler
  ): Observable<HttpEvent<any>> {
    return next.handle(request).pipe(
      catchError((error: HttpErrorResponse) => {
        switch (error.status) {
          case 0:
            this.openDialog(
              "No contact",
              "Unable to establish connection to server"
            );
            break;
          case 400:
            this.openDialog("Bad request", error.error.error.message);
            break;
          case 500:
            this.openDialog("Internal server error", error.error.error.message);
            break;
          case 503:
            this.openDialog("Service unavailable", error.error.error.message);
            break;
          default:
            this.openDialog("Unexpected error", error.statusText);
            break;
        }

        return throwError(error);
      })
    );
  }
  openDialog(title: string, lineOne: string, lineTwo?: string) {
    const dialog = this.dialog.open(ErrorHandlingDialogComponent, {
      // width: "600px",
      data: { title, lineOne, lineTwo },
    });
  }
}
