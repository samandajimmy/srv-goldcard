ALTER TABLE personal_informations
ADD address_province VARCHAR(255),
ADD address_village VARCHAR(255),
ADD address_subdistrict VARCHAR(255);

ALTER TABLE occupations
ADD office_province VARCHAR(255),
ADD office_village VARCHAR(255),
ADD office_subdistrict VARCHAR(255);