import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { AddressModule } from './address/address.module';
import { AppDataSource } from './db/config';

@Module({
  imports: [
    AddressModule,
    TypeOrmModule.forRoot(AppDataSource),
  ],
  controllers: [],
  providers: [],
})
export class AppModule {}
