import { Routes } from '@angular/router';
import { LoginComponent } from './pages/login.component';
import { RegisterComponent } from './pages/register.component';
import { ConcertsComponent } from './pages/concerts.component';
import { SetlistComponent } from './pages/setlist.component';

export const routes: Routes = [
  { path: '', redirectTo: 'concerts', pathMatch: 'full' },
  { path: 'login', component: LoginComponent },
  { path: 'register', component: RegisterComponent },
  { path: 'concerts', component: ConcertsComponent },
  { path: 'concerts/:id/setlist', component: SetlistComponent },
  { path: '**', redirectTo: 'concerts' },
];
