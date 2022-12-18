import { Address } from './address.model';
import { Test } from '@nestjs/testing';
import * as request from 'supertest';
import { HttpStatus, INestApplication, ValidationPipe } from '@nestjs/common';
import { AppTestModule } from 'src/app.test.module';

let mockAddress = function(): Address {
    return {
        id: 0,
        zipCode: "00000000",
        address: "mock address",
        district: "mock district",
        city: "mock city",
        state: "MS"
    };
}

let mockInvalidAddress = function(): Address {
    return {
        id: 0,
        zipCode: "",
        address: "",
        district: "",
        city: "",
        state: ""
    };
}

const addressBadRequestExpectedBody = {
    statusCode: 400,
    message: [
        'zipCode must be longer than or equal to 8 characters',
        'zipCode should not be empty',
        'address should not be empty',
        'district should not be empty',
        'city should not be empty',
        'state should not be empty',
        'state must be longer than or equal to 2 characters'
    ],
    error: 'Bad Request'
}

describe("AddressController", () => {

    let app: INestApplication;

    beforeEach(async () => {
        const addressModule = await Test
            .createTestingModule({imports: [AppTestModule]})
            .compile();
        app = addressModule.createNestApplication();
        app.useGlobalPipes(new ValidationPipe());
        await app.init();
    });

    it("/POST address 201", () => {
        let address = mockAddress()
        return request(app.getHttpServer())
            .post("/address")
            .send(address)
            .expect(HttpStatus.CREATED.valueOf())
    });

    it("/POST address 400", () => {
        let address = mockInvalidAddress()
        return request(app.getHttpServer())
            .post("/address")
            .send(address)
            .expect(addressBadRequestExpectedBody)
            .expect(HttpStatus.BAD_REQUEST.valueOf())
    });
});