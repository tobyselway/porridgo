struct CameraUniform {
    view: mat4x4<f32>,
    projection: mat4x4<f32>,
}

@group(1)
@binding(0)
var<uniform> camera: CameraUniform;

struct InstanceInput {
    @location(5) model_matrix_0: vec4<f32>,
    @location(6) model_matrix_1: vec4<f32>,
    @location(7) model_matrix_2: vec4<f32>,
    @location(8) model_matrix_3: vec4<f32>,
}

struct VertexInput {
    @location(0) position: vec3<f32>,
    @location(1) uv_mapping: vec2<f32>,
    @location(2) normals: vec3<f32>,
}

struct VertexOutput {
    @builtin(position) clip_position: vec4<f32>,
    @location(1) uv_mapping : vec2<f32>,
    @location(2) normals: vec3<f32>,
}

@vertex
fn vs_main(
    model: VertexInput,
    instance: InstanceInput,
) -> VertexOutput {
    let model_matrix = mat4x4<f32>(
        instance.model_matrix_0,
        instance.model_matrix_1,
        instance.model_matrix_2,
        instance.model_matrix_3,
    );

    var out: VertexOutput;
    out.uv_mapping = model.uv_mapping;
    out.clip_position = camera.projection * camera.view * model_matrix * vec4<f32>(model.position, 1.0);
    out.normals = model.normals;
    return out;
}

@group(0)
@binding(0)
var t_diffuse: texture_2d<f32>;

@group(0)
@binding(1)
var s_diffuse: sampler;

@fragment
fn fs_main(vertex : VertexOutput) -> @location(0) vec4<f32> {
    return textureSample(t_diffuse, s_diffuse, vertex.uv_mapping);
}
