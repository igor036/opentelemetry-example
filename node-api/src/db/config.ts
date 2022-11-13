import { Address } from "src/address/address.model";
import { TypeOrmModuleOptions } from '@nestjs/typeorm';

export const AppDataSource: TypeOrmModuleOptions = {
    type: 'postgres',
    host: '127.0.0.1',
    port: 5432,
    username: 'root',
    password: 'root',
    database: 'test',
    entities: [Address],
    synchronize: true,
}